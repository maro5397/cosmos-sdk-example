# cosmos-sdk-example
cosmos sdk를 사용한 로컬 블록체인 네트워크 예제

## introduce
GCP에 Cosmos SDK Validator를 운용하기 위한 실험 환경 구축

## setting
### env 파일 설정
아래 .env 설정을 monitoring 디렉토리 위치에 저장
```ini
GRAFANA_ADMIN_USER = [userId]
GRAFANA_ADMIN_PASSWORD = [password]
```
### cosmos-sdk build
cosmos-sdk를 빌드한 후 simd파일을 /usr/local/bin 위치로 이동 (^v0.53.4)
```shell
$ git clone https://github.com/cosmos/cosmos-sdk
$ cd cosmos-sdk
$ git checkout v0.53.4
$ make build
$ cp ./build/simd /usr/local/bin
```

## how to start
```shell
$ ansible-playbook -i hosts.ini playbook.yml
```