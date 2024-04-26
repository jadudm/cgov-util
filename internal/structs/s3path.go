package structs

type S3Path struct {
	Bucket string
	Key    string
}

type Thunk func()
type Choice struct {
	Local  Thunk
	Remote Thunk
}
