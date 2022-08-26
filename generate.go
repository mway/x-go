package x

//go:generate rm -rf testing/internal/requiremock
//go:generate mockgen -package requiremock -destination testing/internal/requiremock/requiremock.go github.com/stretchr/testify/require TestingT
