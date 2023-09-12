package registry

import (
	"fmt"
	"github.com/opencontainers/go-digest"
	"imageRemover/logger"
)

func getDigest(ServerUrl, user, pass, repo, tag string) (digest.Digest, error) {
	headers := makeDockerHeaders(user, pass)
	tagUrl := fmt.Sprintf("%s/v2/%s/manifests/%s", ServerUrl, repo, tag)
	tagResp, err := doRequest("GET", tagUrl, &headers)
	if err != nil {
		logger.Log.Warn().Err(err).Str("repo", repo).Str("tag", tag).
			Msg("Can't get manifest from registry")
		return "", err
	}

	dig, err := digest.Parse(tagResp.Header.Get("Docker-Content-Digest"))
	if err != nil {
		logger.Log.Warn().Err(err).Str("repo", repo).Str("tag", tag).
			Msg("Can't parse digest from registry")
		return "", err
	}
	if err := tagResp.Body.Close(); err != nil {
		logger.Log.Error().Err(err).Msg("Can't close tag response")
	}
	return dig, nil
}
