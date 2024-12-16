pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                // Checkout the code from the Git repository
                git url: 'https://github.com/jamil225/chat-app.git', branch: 'main'
            }
        }

        stage('Build') {
            steps {
                // Navigate to the chat-service directory
                dir('chat-service') {
                    // Set up Go environment
                    sh 'export GOPATH=$HOME/go'
                    sh 'export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin'

                    // Install dependencies
                    sh 'go mod tidy'

                    // Build the project
                    sh 'go build -o chat-service'
                }
            }
        }

        stage('Test') {
            steps {
                // Navigate to the chat-service directory
                dir('chat-service') {
                    // Run tests
                    sh 'go test ./...'
                }
            }
        }

        stage('Run') {
            steps {
                // Navigate to the chat-service directory
                dir('chat-service') {
                    // Run the application
                    sh './chat-service'
                }
            }
        }
    }

    post {
        always {
            // Clean up workspace
            cleanWs()
        }
    }
}
