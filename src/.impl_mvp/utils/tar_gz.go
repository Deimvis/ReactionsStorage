package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"reflect"

	"gopkg.in/yaml.v3"
)

func CreateTarGz(in interface{}) ([]byte, error) {
	var tarGzContent bytes.Buffer

	tarWriter := tar.NewWriter(&tarGzContent)
	defer tarWriter.Close()

	v := reflect.ValueOf(in)
	for i := 0; i < v.NumField(); i++ {
		fileName := v.Type().Field(i).Tag.Get("filename")
		content, err := yaml.Marshal(v.Field(i).Interface())
		if err != nil {
			return nil, err
		}
		header := &tar.Header{
			Name: fileName,
			Size: int64(len(content)),
		}

		err = tarWriter.WriteHeader(header)
		if err != nil {
			return nil, err
		}
		_, err = tarWriter.Write(content)
		if err != nil {
			return nil, err
		}
	}

	var gzContent bytes.Buffer
	gzWriter := gzip.NewWriter(&gzContent)

	if _, err := io.Copy(gzWriter, &tarGzContent); err != nil {
		return nil, err
	}

	if err := gzWriter.Close(); err != nil {
		return nil, err
	}

	return gzContent.Bytes(), nil
}

func ExtractTarGz(in []byte, out interface{}) error {
	gzReader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return err
	}
	defer gzReader.Close()
	tarReader := tar.NewReader(gzReader)

	filename2field := getTag2FieldAddr(out, "filename")
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}
		content, err := io.ReadAll(tarReader)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(content, filename2field[header.Name])
		if err != nil {
			return err
		}
	}
	return nil
}

func getTag2FieldAddr(obj interface{}, tag string) map[string]interface{} {
	res := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(tag)
		res[tag] = v.Field(i).Addr().Interface()
	}
	return res
}
