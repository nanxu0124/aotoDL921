FROM nvidia/cuda:12.2.0-devel-ubuntu20.04

ENV PIP_ROOT_USER_ACTION=ignore
ENV LANG="C.UTF-8"

WORKDIR /root

# own config
COPY myconfig myconfig
RUN mv myconfig/condarc .condarc \
	&& mv myconfig/pydistutils.cfg .pydistutils.cfg \
	&& mv myconfig/Shanghai /etc/localtime \
	&& mkdir -p .ssh \
	&& mv myconfig/trans-docker.pub .ssh/ 

# Avoid Public GPG key error
# https://github.com/NVIDIA/nvidia-docker/issues/1631
RUN rm /etc/apt/sources.list.d/cuda.list 

# use Mirror to speed up downloads
RUN sed -i 's/http:\/\/archive.ubuntu.com\/ubuntu\//http:\/\/mirrors.aliyun.com\/ubuntu\//g' /etc/apt/sources.list

# Install the required packages
RUN apt update && apt install -y build-essential apt-utils ca-certificates vim wget ssh \
    && apt-get update \
    && apt-get install -y ffmpeg libsm6 libxext6 git ninja-build libglib2.0-0 libsm6 libxrender-dev libxext6 net-tools\
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# set ssh key
RUN cat .ssh/trans-docker.pub >> .ssh/authorized_keys \
        && echo "/etc/init.d/ssh start" >> ~/.bashrc

EXPOSE 22

# ssh login
RUN cd /etc/ssh \
        && echo "PermitRootLogin yes" >> sshd_config 
        #&& sed -i "/PubkeyAuthentication/s/^#//" sshd_config \
        #&& sed -i "/AuthorizedKeysFile/s/^#//" sshd_config

# Install Miniconda
RUN wget -q -O Miniconda.sh https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh \
        && /bin/bash Miniconda.sh -b -p miniconda
ENV PATH /root/miniconda/bin:$PATH
RUN /bin/bash -c "source miniconda/bin/activate" && conda init

# Install python and torch
RUN conda create -y -n workspace python=3.8
RUN echo "source activate workspace" >> ~/.bashrc
ENV PATH /root/miniconda/envs/workspace/bin:$PATH
RUN pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple \
	&& pip install --upgrade pip
RUN pip3 install torch torchvision torchaudio


ENTRYPOINT bash
