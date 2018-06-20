package scm

import (
	"context"
	"errors"
	"strconv"
)

// FakeSCM implements the SCM interface.
type FakeSCM struct {
	Repositories map[uint64]*Repository
	Directories  map[uint64]*Directory
	Hooks        map[uint64]int
}

// NewFakeSCMClient returns a new Fake client implementing the SCM interface.
func NewFakeSCMClient() *FakeSCM {
	return &FakeSCM{
		Repositories: make(map[uint64]*Repository),
		Directories:  make(map[uint64]*Directory),
		Hooks:        make(map[uint64]int),
	}
}

// ListDirectories implements the SCM interface.
func (s *FakeSCM) ListDirectories(ctx context.Context) ([]*Directory, error) {
	var dirs []*Directory
	for _, dir := range s.Directories {
		dirs = append(dirs, dir)
	}

	return dirs, nil
}

// CreateDirectory implements the SCM interface.
func (s *FakeSCM) CreateDirectory(ctx context.Context, opt *CreateDirectoryOptions) (*Directory, error) {
	id := len(s.Directories) + 1
	dir := &Directory{
		ID:     uint64(id),
		Path:   opt.Path,
		Avatar: "https://avatars3.githubusercontent.com/u/1000" + strconv.Itoa(id) + "?v=3",
	}
	s.Directories[dir.ID] = dir
	return dir, nil
}

// GetDirectory implements the SCM interface.
func (s *FakeSCM) GetDirectory(ctx context.Context, id uint64) (*Directory, error) {
	dir, ok := s.Directories[id]
	if !ok {
		return nil, errors.New("directory not found")
	}
	return dir, nil
}

// CreateRepository implements the SCM interface.
func (s *FakeSCM) CreateRepository(ctx context.Context, opt *CreateRepositoryOptions) (*Repository, error) {
	id := len(s.Repositories) + 1
	repo := &Repository{
		ID:          uint64(id),
		Path:        opt.Path,
		WebURL:      "https://example.com/" + opt.Directory.Path + "/" + opt.Path,
		SSHURL:      "git@example.com:" + opt.Directory.Path + "/" + opt.Path,
		HTTPURL:     "https://example.com/" + opt.Directory.Path + "/" + opt.Path + ".git",
		DirectoryID: opt.Directory.ID,
	}
	s.Repositories[repo.ID] = repo
	return repo, nil
}

// GetRepositories implements the SCM interface.
func (s *FakeSCM) GetRepositories(ctx context.Context, directory *Directory) ([]*Repository, error) {
	var repos []*Repository
	for _, repo := range s.Repositories {
		if repo.DirectoryID == directory.ID {
			repos = append(repos, repo)
		}
	}
	return repos, nil
}

// DeleteRepository implements the SCM interface.
func (s *FakeSCM) DeleteRepository(ctx context.Context, id uint64) error {
	if _, ok := s.Repositories[id]; !ok {
		return errors.New("repository not found")
	}
	delete(s.Repositories, id)
	return nil
}

// ListHooks implements the SCM interface.
func (s *FakeSCM) ListHooks(ctx context.Context, repo *Repository) ([]*Hook, error) {
	// TODO no implementation provided yet
	return nil, nil
}

// CreateHook implements the SCM interface.
func (s *FakeSCM) CreateHook(ctx context.Context, opt *CreateHookOptions) error {
	if _, ok := s.Repositories[opt.Repository.ID]; !ok {
		return errors.New("repository not found")
	}
	s.Hooks[opt.Repository.ID]++
	return nil
}

// CreateTeam implements the SCM interface.
func (s *FakeSCM) CreateTeam(ctx context.Context, opt *CreateTeamOptions) (*Team, error) {
	// TODO no implementation provided yet
	return nil, nil
}

// CreateCloneURL implements the SCM interface.
func (s *FakeSCM) CreateCloneURL(ctx context.Context, opt *CreateClonePathOptions) (string, error) {
	return "", nil
}

// AddTeamRepo implements the SCM interface.
func (s *FakeSCM) AddTeamRepo(ctx context.Context, opt *AddTeamRepoOptions) error {
	return nil
}

// GetUserNameByID implements the SCM interface.
func (s *FakeSCM) GetUserNameByID(ctx context.Context, remoteID uint64) (string, error) {
	return "", nil
}
