# LogAid Testing Environment
FROM ubuntu:22.04

# Install common tools that LogAid helps with
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    git \
    vim \
    nano \
    htop \
    tree \
    zip \
    unzip \
    build-essential \
    python3 \
    python3-pip \
    nodejs \
    npm \
    golang-go \
    && rm -rf /var/lib/apt/lists/*

# Install Docker (for testing Docker plugin)
RUN curl -fsSL https://get.docker.com -o get-docker.sh && \
    sh get-docker.sh && \
    rm get-docker.sh

# Install kubectl (for testing Kubernetes plugin)
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
    rm kubectl

# Create a test user
RUN useradd -m -s /bin/bash testuser && \
    usermod -aG docker testuser

# Copy LogAid binary (will be built separately)
COPY logaid /usr/local/bin/logaid
RUN chmod +x /usr/local/bin/logaid

# Set up test environment
USER testuser
WORKDIR /home/testuser

# Create some test scenarios
RUN echo '#!/bin/bash' > test-scenarios.sh && \
    echo 'echo "Testing LogAid with various error scenarios..."' >> test-scenarios.sh && \
    echo 'echo "1. Testing apt typo..."' >> test-scenarios.sh && \
    echo 'apt instal curl' >> test-scenarios.sh && \
    echo 'echo "2. Testing git typo..."' >> test-scenarios.sh && \
    echo 'git checout main' >> test-scenarios.sh && \
    echo 'echo "3. Testing npm typo..."' >> test-scenarios.sh && \
    echo 'npm instal express' >> test-scenarios.sh && \
    echo 'echo "4. Testing docker typo..."' >> test-scenarios.sh && \
    echo 'docker rnu hello-world' >> test-scenarios.sh && \
    echo 'echo "5. Testing kubectl typo..."' >> test-scenarios.sh && \
    echo 'kubectl gt pods' >> test-scenarios.sh && \
    chmod +x test-scenarios.sh

CMD ["bash"]
