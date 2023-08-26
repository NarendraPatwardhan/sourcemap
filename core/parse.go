package core

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"

	"machinelearning.one/sourcemap/compose/logger"
)

func Parse(ctx context.Context, addr string, opts ParseOpts) Repository {
	lg := logger.Get(ctx)

	r, err := git.CloneContext(ctx, memory.NewStorage(), nil, &git.CloneOptions{
		URL:   addr,
		Depth: opts.Limit,
	})
	if err != nil {
		lg.Fatal().Err(err).Msg("failed to clone repository")
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

	for len(queue) > 0 {
		co := queue[0]
		queue = queue[1:]

		commit := &Commit{
			Hash:    co.Hash.String(),
			Author:  co.Author.Name,
			Message: co.Message,
			Time:    co.Author.When.String(),
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
		finalize(ctx, data)

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

	for _, entry := range tree.Entries {
		path := path + "/" + entry.Name
		if filter(path, excludeGlobs, excludePaths) {
			continue
		}

		child := &Data{
			Name: entry.Name,
			Path: path,
		}

		isFile := entry.Mode.IsFile()
		if isFile {

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

			extract(ctx, child, tree, path, excludeGlobs, excludePaths)
		}

		data.Children = append(data.Children, child)
	}
}

func filter(path string, excludeGlobs, excludePaths []string) bool {
	return false
}

func finalize(ctx context.Context, data *Data) {
	for _, child := range data.Children {
		finalize(ctx, child)
		data.Size += child.Size
	}
}
