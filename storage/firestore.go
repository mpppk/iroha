package storage

import (
	"context"
	"log"

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
	WordsList []*fireStoreWords
}

type fireStoreWords struct {
	Id    int
	Words []*fireStoreWord
}

type fireStoreWord struct {
	Id   int
	Bits int
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
	fswordsList := fromWordsList(wordsList)
	doc := &cacheDoc{WordsList: fswordsList}
	log.Println("doc", len(doc.WordsList))
	_, err := f.client.Collection(f.collectionName).Doc(toStorageStrKey(indices)).Set(ctx, doc)
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
	//var WordsList [][]*fireStoreWord
	//var WordsList [][]*ktkn.Word
	var doc cacheDoc
	if err := dsnap.DataTo(&doc); err != nil {
		return nil, false, errors.Wrap(err, "failed to convert firestore data to WordsList")
	}
	wordsList := toWordsList(doc.WordsList)
	log.Println("firestore get")
	log.Println(doc.WordsList)

	return wordsList, true, nil
}

func fromWord(word *ktkn.Word) *fireStoreWord {
	return &fireStoreWord{
		Id:   int(word.Id),
		Bits: int(word.Bits),
	}
}

func fromWords(id int, words []*ktkn.Word) *fireStoreWords {
	fswords := &fireStoreWords{Id: id, Words: []*fireStoreWord{}}
	for _, word := range words {
		fswords.Words = append(fswords.Words, fromWord(word))
	}
	return fswords
}

func fromWordsList(wordsList [][]*ktkn.Word) (fswordsList []*fireStoreWords) {
	for i, words := range wordsList {
		fswordsList = append(fswordsList, fromWords(i, words))
	}
	return
}

func toWord(fsword *fireStoreWord) *ktkn.Word {
	return &ktkn.Word{
		Id:   ktkn.WordId(fsword.Id),
		Bits: ktkn.WordBits(fsword.Bits),
	}
}

func toWords(fswords *fireStoreWords) (words []*ktkn.Word) {
	for _, fsword := range fswords.Words {
		words = append(words, toWord(fsword))
	}
	return
}

func toWordsList(fswordsList []*fireStoreWords) (wordsList [][]*ktkn.Word) {
	for _, fswords := range fswordsList {
		wordsList = append(wordsList, toWords(fswords))
	}
	return
}
