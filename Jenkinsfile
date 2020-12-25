pipeline {
    agent any 
    stages {
        stage('Stage 1') {
            steps {
                echo 'Hello world!'
                sh 'docker ps'
                sh 'echo 'Hello world!' > hw.txt'
            }
        }
    }
}
