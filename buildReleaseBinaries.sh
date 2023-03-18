#! /bin/bash
binaryname=$(cat go.mod | grep module | head -n 1 | cut -d" " -f2 | rev | cut -d/ -f1 | rev)

rm -rf uploader-binaries
mkdir uploader-binaries

# building
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -extldflags=-static" -gcflags=-trimpath=x/y -o uploader-binaries/linux-amd64/
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -extldflags=-static" -gcflags=-trimpath=x/y -o uploader-binaries/mac-Intel-amd64/
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w -extldflags=-static" -gcflags=-trimpath=x/y -o uploader-binaries/mac-M1-arm64/
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -extldflags=-static" -gcflags=-trimpath=x/y -o uploader-binaries/windows-amd64/

cd uploader-binaries

# visual verifying
for i in *
do
cd $i
echo $i
file *
echo
cd ..
done


# zipping
zip ../binaries-$binaryname-vREPLACE-VERSION-HERE.zip -r *
cd ..
rm -rf uploader-binaries