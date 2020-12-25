pipeline {
    agent any 
    stages {
        stage('Hello world console') {
            steps {
                echo 'Hello world!'
                sh 'docker ps'
                sh 'echo "Hello world!" > hw.txt'
            }
        }
        
        stage('Test docker') {
            steps {
                sh 'docker ps'
            }
        }
        
        stage('Save artifact') {
            steps {
                sh 'echo "Hello world!" > hw.txt'
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: 'hw.txt', onlyIfSuccessful: true
        }
    }
}
