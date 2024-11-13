# lazy-mzstudy-go
A Solution for Lazzzzzzzzy person 😒


## **주의 사항** ⚠️
- 🚨 **(필수 확인) 절대 메가존 런쳐(팝스)의 아이디와 패스워드를 사용하지 마세요!**   
- 🚨 **(필수 확인) 법정의무교육 사이트의 아이디와 패스워드를 사용해야 합니다.**
- **프로그램 실행 도중 `https://mz.livestudy.com`에 로그인하지 마세요!**   
    *중복 세션 방지 기능으로 인해 기존 로그인이 해제될 수 있습니다.*
- 서버 또는 프로그램 로직 상의 이슈로 인해 한 번의 실행으로 모든 학습이 완료되지 않을 수 있습니다. 따라서 여러 번 실행해야 할 수 있습니다.


## 사용 방법

## 1단계: [다운로드](https://github.com/kktae/lazy-mzstudy-go/releases)
1. [릴리즈 페이지](https://github.com/kktae/lazy-mzstudy-go/releases)에 접속하세요.
2. 운영체제(OS)와 아키텍처에 맞는 파일을 찾아 다운로드하세요.   

    > 제공되는 운영체제는 `mac`, `linux`, `windows`가 있고,   
    > 제공되는 아키텍처는 `x86_64(amd64)`, `arm64`가 있습니다.  

    > 파일명 형식은 `lazy-mzstudy-go_<YOUR_OS>_<YOUR_ARCH>.tar.gz` 입니다.   

    - 예시)   
        - `macOS Apple Silicon`: lazy-mzstudy-go_**Darwin**_**arm64**.tar.gz   
        - `Windows 64비트`: lazy-mzstudy-go_**Windows**_**x86_64**.zip

## 2단계: 압축 해제
- 다운로드한 파일의 압축을 해제합니다.
    - 예시: `tar -xzvf lazy-mzstudy-go_Darwin_arm64.tar.gz`

## 3단계: 설정 파일 생성
1. 압축을 해제한 폴더에 `config.yaml` 파일을 생성합니다.
1. `config.yaml` 파일에 나의 mzstudy 아이디와 비밀번호를 입력합니다.
    ```yaml
    # config.yaml
    mz_id: "your_account@mz.co.kr"
    mz_pw: "your_password"
    ```
## 4단계: 실행
1. 실행 파일을 실행합니다.
1. 프로그램 실행 중 출력되는 로그를 확인합니다.
    - `skip already study lesson...`: **이미 완료된 학습**입니다.
    - `begin study lesson...`: **아직 완료되지 않은 학습**입니다.
1. 실행이 완료될 때까지 기다립니다. (약 1분 30초 소요)

## 5단계: 학습 완료 확인
1. 프로그램을 다시 실행합니다.
1. 모든 학습 항목에 대해 `skip already study lesson...` 메시지가 출력되는지 확인합니다.
1. 만약 완료되지 않은 학습이 있다면, **`4단계: 실행`** 부터 다시 반복합니다.
1. 모든 학습이 완료되었다면, https://mz.livestudy.com 에 로그인하여 학습 완료 상태를 최종 확인합니다.
