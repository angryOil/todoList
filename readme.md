# 필독 브랜치 규칙
main - prod, stage - stage , dev - 개발서버

pr 순서 feat => dev => (hotfix/bug)stage => main


1. 최신 dev브랜치에서 feature 만들기
2. dev에 push전 dev pull 받기 
3. bug/hotfix 를 제외한 브랜치(ex:feat)로 main/stage에 직접pr금지 


## todoList
할일 목록을 보여주는 간단한 api 입니다.

### makefile

```shell
# local postgres run (docker-compose)
make local-db
# local postgres migrate init
make local-init
# local postgres apply migrate
make local-migrate
```
