package storage

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/mpppk/iroha/ktkn"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type FireStore struct {
	client         *firestore.Client
	collectionName string
}

type Cache struct {
	indices   []int
	wordsList [][]*ktkn.Word
}

type cacheDoc struct {
	WordsList string
}

func NewFireStore(ctx context.Context, filePath string) (storage *FireStore, err error) {
	var baseErrMsg = "failed to create firestore storage: "
	sa := option.WithCredentialsFile(filePath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, errors.Wrap(err, baseErrMsg+"failed to create firebase new app")
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, errors.Wrap(err, baseErrMsg+"failed to create firestore client")
	}

	return &FireStore{
		client:         client,
		collectionName: "cache",
	}, errors.Wrap(err, baseErrMsg+"failed to close firestore client")
}

func (f *FireStore) Set(indices []int, wordsList [][]*ktkn.Word) error {
	ctx := context.Background()
	wl := wordsList
	if wl == nil {
		wl = make([][]*ktkn.Word, 0)
	}
	wordsListJsonBytes, err := json.Marshal(wordsList)
	doc := &cacheDoc{WordsList: string(wordsListJsonBytes)}
	_, err = f.client.Collection(f.collectionName).Doc(toStorageStrKey(indices)).Set(ctx, doc)
	if err != nil {
		return errors.Wrap(err, "failed to set cache to firestore")
	}
	return nil
}

func (f *FireStore) Get(indices []int) ([][]*ktkn.Word, bool, error) {
	ctx := context.Background()
	key := toStorageStrKey(indices)
	dsnap, err := f.client.Collection(f.collectionName).Doc(key).Get(ctx)
	if dsnap != nil && !dsnap.Exists() {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, errors.Wrapf(err, "failed to get cache from firestore. indices: %s", key)
	}
	var doc cacheDoc
	if err := dsnap.DataTo(&doc); err != nil {
		return nil, false, errors.Wrap(err, "failed to convert firestore data to WordsList")
	}

	var wordsList [][]*ktkn.Word
	if err = json.Unmarshal([]byte(doc.WordsList), &wordsList); err != nil {
		return nil, false, errors.Wrap(err, "failed to unmarshal firestore results")
	}
	return wordsList, true, nil
}
