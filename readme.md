# 필독 브랜치 규칙
main - prod, stage - stage , dev - 개발서버

1. 최신 dev브랜치에서 feature 만들기
2. push전 dev pull 받기 
3. bug/hotfix 를 제외한 브랜치(ex:feat) main/stage로 pr금지 

### todoList

### makefile

```shell
# local postgres run (docker-compose)
make local-db
# local postgres migrate init
make local-init
# local postgres apply migrate
make local-migrate
```
