# lazy-mzstudy-go
to lazzzzzzy

## how to use
1. 운영체제에 맞는 바이너리 다운받기
2. `config.yaml` 생성하기   
2.1. 아이디 및 패스워드 입력하기
    ```yaml
    # config.yaml
    mz_id: "your_account@mz.co.kr"
    mz_pw: "your_password"
    ```
4. 바이너리 실행하기   
4.1. 완료된 학습의 경우 `skip already study lesson...`   
4.2. 완료되지 않은 학습의 경우 `begin study lesson...`
5. 실행이 완료될 때까지 대기(대략 2분)
6. 바이너리 재실행 후, **4.2번**과 완료되지 않은 학습이 있는 지 확인
7. 이하 반복, 모든 학습이 완료되면 **4.1번**과 같이 전부 skip 됨
8. https://mz.livestudy.com 로그인 후, 학습 완료 상태 확인

## 주의사항
* 프로세스가 실행되는 동안에는 절대 mzstudy 로그인하면 안 됨   
    *로그인 시, 중복 세션 방지로 인해 기존 로그인이 해체됨*