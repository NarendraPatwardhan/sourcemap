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

		repo = append(repo, commit)
		co.Parents().ForEach(func(parent *object.Commit) error {
			queue = append(queue, parent)
			return nil
		})
	}

	return repo
}
