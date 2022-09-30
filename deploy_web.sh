# set -x
git pull
go build
nohup ./fcc >fcc.log 2>&1 &
