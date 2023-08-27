package core

import (
	"context"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"

	"machinelearning.one/sourcemap/compose/logger"
	"machinelearning.one/sourcemap/core/ext"
)

func Parse(ctx context.Context, addr string, opts ParseOpts) Repository {
	lg := logger.Get(ctx)

	depth := opts.Limit
	if depth != 0 {
		depth += 1 // Fetch an extra commit to obtain changes
	}

	r, err := git.CloneContext(ctx, memory.NewStorage(), nil, &git.CloneOptions{
		URL:   addr,
		Depth: depth,
	})
	if err != nil {
		lg.Fatal().Err(err).Msgf("failed to clone repository at %s", addr)
	}

	ref, err := r.Head()
	if err != nil {
		lg.Fatal().Err(err).Msg("failed to get head")
	}

	co, err := r.CommitObject(ref.Hash())
	if err != nil {
		lg.Fatal().Err(err).Msg("failed to get commit")
	}

	repo := Repository{}

	queue := make([]*object.Commit, 0)
	queue = append(queue, co)
	seen := 0

	for len(queue) > 0 {
		co := queue[0]
		queue = queue[1:]
		seen += 1

		if opts.Limit != 0 && seen > opts.Limit {
			break
		}

		stats, err := co.Stats()
		if err != nil {
			lg.Warn().
				Err(err).
				Msgf("failed to get stats for %s", co.Hash.String())
		}
		statsLookup := make(map[string]Changes, len(stats))
		for _, stat := range stats {
			statsLookup[stat.Name] = Changes{
				Addition: stat.Addition,
				Deletion: stat.Deletion,
			}
		}

		commit := &Commit{
			Hash:      co.Hash.String(),
			Author:    co.Author.Name,
			Message:   co.Message,
			Timestamp: co.Author.When.String(),
		}

		data := &Data{
			Name:     "<root>",
			Path:     "",
			Children: make([]*Data, 0),
		}

		tree, err := co.Tree()
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to get root tree")
		}

		extract(ctx, data, tree, data.Path, opts.ExcludeGlobs, opts.ExcludePaths)
		gather(ctx, data, statsLookup)

		commit.Data = data

		repo = append(repo, commit)
		co.Parents().ForEach(func(parent *object.Commit) error {
			queue = append(queue, parent)
			return nil
		})
	}

	return repo
}

func extract(
	ctx context.Context,
	data *Data,
	tree *object.Tree,
	path string,
	excludeGlobs, excludePaths []string,
) {
	lg := logger.Get(ctx)

	if path != "" {
		path += "/"
	}

	for _, entry := range tree.Entries {

		if filter(ctx, path+entry.Name, excludeGlobs, excludePaths) {
			continue
		}

		child := &Data{
			Name: entry.Name,
			Path: path + entry.Name,
		}

		isFile := entry.Mode.IsFile()
		if isFile {

			child.Repr = ext.Get(entry.Name)

			blob, err := tree.TreeEntryFile(&entry)
			if err != nil {
				lg.Warn().Err(err).Msgf("failed to get blob for %s", path)
				continue
			}

			child.Size = blob.Size
		} else {

			child.Children = make([]*Data, 0)
			tree, err := tree.Tree(entry.Name)
			if err != nil {
				lg.Warn().Err(err).Msgf("failed to get tree for %s", path)
				continue
			}

			extract(ctx, child, tree, path+entry.Name, excludeGlobs, excludePaths)
		}

		data.Children = append(data.Children, child)
	}
}

func filter(ctx context.Context, path string, excludeGlobs, excludePaths []string) bool {
	lg := logger.Get(ctx)

	if len(excludeGlobs) > 0 {
		for _, glob := range excludeGlobs {
			match, err := doublestar.Match(glob, path)
			if err != nil {
				lg.Warn().
					Err(err).
					Msgf("failed to use glob %s as exclusion criteria, pattern malformed", glob)
				continue
			}
			if match {
				return true
			}
		}
	}

	if len(excludePaths) > 0 {
		for _, excludePath := range excludePaths {
			if strings.Contains(path, excludePath) {
				return true
			}
		}
	}

	return false
}

func gather(ctx context.Context, data *Data, stats map[string]Changes) {
	for _, child := range data.Children {
		gather(ctx, child, stats)
		data.Size += child.Size
		changes, ok := stats[child.Path]
		if ok {
			child.Changes.Addition = changes.Addition
			child.Changes.Deletion = changes.Deletion
			data.Changes.Addition += changes.Addition
			data.Changes.Deletion += changes.Deletion
		}
	}
}
