package helpers

import (
	"github.com/cascad/goloot/data_structs"
	"github.com/cascad/goloot/erlang"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/bkaradzic/go-lz4"
	"github.com/shamaton/msgpack"
	"github.com/zegl/goriak/v3"
)

func GetDataFromRiak(uid string, conn *goriak.Session, bucket string) (*[]byte, error) {
	btype := "default"

	var rawInfo []byte

	_, err := goriak.Bucket(bucket, btype).GetRaw(uid, &rawInfo).Run(conn)
	if err != nil {
		return &rawInfo, err
	}

	return &rawInfo, nil
}

func DecodeDBData(uid string, raw *[]byte) (*[]byte, error) {
	//raw, err := ioutil.ReadFile("/home/cascad/afterget.bin")
	//parsers.PanicOnErr(err)
	term, err := erlang.BinaryToTerm(*raw)
	//parsers.PanicOnErr(err)
	if err != nil {
		return raw, err
	}

	tup, ok := term.(erlang.OtpErlangTuple)
	//parsers.PanicOnNok(ok)
	if !ok {
		return raw, fmt.Errorf("%s -> %s", "bad cast to OtpErlangTuple", uid)
	}
	bin, ok := tup[1].(erlang.OtpErlangBinary)
	//parsers.PanicOnNok(ok)
	if !ok {
		return raw, fmt.Errorf("%s -> %s", "bad cast to OtpErlangBinary", uid)
	}
	//parsers.PanicOnErr(ioutil.WriteFile("/home/cascad/beforeDec.bin", bin.Value, 0644))

	decoded := make([]byte, base64.StdEncoding.EncodedLen(len(bin.Value)))
	n, err := base64.StdEncoding.Decode(decoded, bin.Value)
	//parsers.PanicOnErr(err)
	if err != nil {
		return raw, fmt.Errorf("%s -> %s", "bad b64 decode", uid)
	}
	decoded = decoded[:n]
	decoded = bytes.TrimLeft(decoded, "\x00")

	//parsers.PanicOnErr(ioutil.WriteFile("/home/cascad/b64dec.bin", decoded, 0644))
	//parsers.PanicOnErr(err)

	ln := len(decoded)
	var swapped []byte
	//swapped := make([]byte, ln, ln)

	if ln < 4 {
		return &decoded, fmt.Errorf("bad data for swap: %s -> %s", uid, decoded)
	}

	swapped = append(swapped, decoded[ln-4:]...)
	swapped = append(swapped, decoded[:len(decoded)-4]...)
	//swapped = bytes.TrimLeft(swapped, "\x00")
	//log.Println("sw", len(swapped))

	//parsers.PanicOnErr(ioutil.WriteFile("/home/cascad/afterDec.bin", swapped, 0644))
	//parsers.PanicOnErr(err)

	var lz4dec []byte //make([]byte, (ln<<8)-ln-2526)
	res, err := lz4.Decode(lz4dec, swapped)
	if err != nil {
		return raw, fmt.Errorf("%s -> %s", "bad lz4 decode", uid)
	}
	//log.Println(string(res))

	//parsers.PanicOnErr(ioutil.WriteFile("/home/cascad/yeah?.txt", res, 0777))

	return &res, nil
}

func GetCoinsFromSyncSave(raw *[]byte) (int, error) {
	term, err := erlang.BinaryToTerm(*raw)
	if err != nil {
		return 0, err
	}

	tup, ok := term.(erlang.OtpErlangTuple)
	//parsers.PanicOnNok(ok)
	if !ok {
		return 0, fmt.Errorf("%s -> %s", "bad cast to OtpErlangTuple")
	}
	bin, ok := tup[1].(erlang.OtpErlangBinary)
	//parsers.PanicOnNok(ok)
	if !ok {
		return 0, fmt.Errorf("%s -> %s", "bad cast to OtpErlangBinary")
	}

	var item data_structs.SyncSave
	err = msgpack.Decode(bin.Value, &item)
	if err != nil {
		return 0, err
	}

	result, err := item.Value()

	return result, nil
}
