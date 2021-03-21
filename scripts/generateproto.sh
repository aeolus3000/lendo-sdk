
protoc -I=./banking/dnb --go_out=. ./banking/dnb/*.proto
protoc -I=./banking --go_out=. ./banking/*.proto
