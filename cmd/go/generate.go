//go:build generate

// //go:generate env -C ../../go/go/src GOOS=darwin GOARCH=amd64 ./make.bash -distpack
// //go:generate env -C ../../go/go/src GOOS=darwin GOARCH=arm64 ./make.bash -distpack
// //go:generate env -C ../../go/go/src GOOS=linux GOARCH=386 ./make.bash -distpack
//go:generate env -C ../../go/go/src GOOS=linux GOARCH=amd64 ./make.bash -distpack
// //go:generate env -C ../../go/go/src GOOS=linux GOARCH=arm64 ./make.bash -distpack
// //go:generate env -C ../../go/go/src GOOS=linux GOARCH=arm ./make.bash -distpack
// //go:generate env -C ../../go/go/src GOOS=windows GOARCH=386 ./make.bash -distpack
// //go:generate env -C ../../go/go/src GOOS=windows GOARCH=amd64 ./make.bash -distpack

package main
