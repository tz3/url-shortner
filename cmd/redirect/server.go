package redirect

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tz3/url-shortner/config"
	"github.com/tz3/url-shortner/storage"
	"github.com/tz3/url-shortner/storage/memory"
	"github.com/tz3/url-shortner/storage/yaml"
	"net/http"
	"net/url"
	"os"
)

const (
	flagHashMapStorage = "with-Hashmap-storage"
	flagYamlStorage    = "with-Yaml-storage"
)

var (
	errorUnSuppoertedStorage = errors.New("ERROR_UNSUPPORTED_STORAGE_TYPE")
	failedStorageSetup       = errors.New("FAILED_STORAGE_SETUP")
)

var (
	serverFlagSet = &pflag.FlagSet{}
)

var strFlags = []string{flagHashMapStorage, flagYamlStorage}

// Serve starts the HTTP server which will do the redirect.
// shortUrl redirected to -> LongUrl (Original Destination)
var Serve = &cobra.Command{
	Use:   "serve",
	Short: "Start the server that handles redirects",
	RunE:  RunServe,
}

func RunServe(cmd *cobra.Command, args []string) error {
	repo, err := getRepository(cmd.Flags())
	if err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reqUrl := &url.URL{
			// There's no support for anything else at this time
			Scheme: "http",
			Host:   r.Host,
			Path:   r.URL.Path,
		}
		ret, err := repo.Get(reqUrl)

		// Iterate though the potential failure modes.
		if errors.Is(err, storage.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", ret.String())
		w.WriteHeader(http.StatusTemporaryRedirect)

	})

	return http.ListenAndServe("localhost:8080", http.DefaultServeMux)
}

func init() {

	// Allow providing the in-memory based storage engine
	serverFlagSet.BoolP(flagHashMapStorage, "m", false, "Use in-memory hashmap storage")
	serverFlagSet.Lookup(flagHashMapStorage).NoOptDefVal = "true"
	err := viper.BindPFlag(config.StorageHashMap, serverFlagSet.Lookup(flagHashMapStorage))
	if err != nil {
		return
	}

	// Allow providing the in-memory based storage engine
	serverFlagSet.StringP(flagYamlStorage, "y", "", "Use the supplied source file as a 'yaml storage'")
	err = viper.BindPFlag(config.StorageYamlFile, serverFlagSet.Lookup(flagYamlStorage))
	if err != nil {
		return
	}

	Serve.Flags().AddFlagSet(serverFlagSet)
	Serve.MarkFlagsOneRequired(strFlags...)
	Serve.MarkFlagsMutuallyExclusive(strFlags...)
}

func getRepository(flags *pflag.FlagSet) (storage.Repository, error) {
	var repo storage.Repository
	for _, f := range strFlags {
		if !flags.Lookup(f).Changed {
			continue
		}

		switch f {
		case flagHashMapStorage:
			// It is possible, in principle, for the user to supply the flag but to be disabling the
			// option rather than just including it.
			if !viper.GetBool(config.StorageHashMap) {
				continue
			}

			return memory.NewHashMap(), nil
		case flagYamlStorage:
			f, err := os.Open(viper.GetString(config.StorageYamlFile))
			if err != nil {
				return nil, fmt.Errorf("%w: %s", failedStorageSetup, err)
			}

			y, err := yaml.New(memory.NewHashMap(), f)
			if err != nil {
				return nil, fmt.Errorf("%w: %s", failedStorageSetup, err)
			}

			return y, nil
		default:
			return nil, errorUnSuppoertedStorage
		}
	}

	if repo == nil {
		return nil, failedStorageSetup
	}

	return repo, nil
}
