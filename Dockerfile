FROM ubuntu:22.04

RUN apt-get update \
 && apt-get install -y --no-install-recommends apt-transport-https curl ca-certificates \
 && apt-get clean \
 && apt-get autoremove \
 && rm -rf /var/lib/apt/lists/* 
RUN mkdir -p /opt/pki/server
RUN mkdir /opt/pki/configs

COPY build/create_aws_config.sh .

ARG CHECK_SUM=bc625c319d96f71c05d899eab04402dc63a455656d46e513b1ea6089b65110ce
RUN curl https://rolesanywhere.amazonaws.com/releases/1.0.4/X86_64/Linux/aws_signing_helper --output /usr/local/bin/aws_signing_helper
RUN echo "${CHECK_SUM}  /usr/local/bin/aws_signing_helper" | sha256sum -c - | grep OK
RUN chmod +x /usr/local/bin/aws_signing_helper

ARG AWS_TRUST_ANCHOR_ARN
ARG AWS_PROFILE_ARN
ARG AWS_ROLE_ARN
ARG AWS_REGION

RUN chmod +x create_aws_config.sh
CMD ["/bin/bash", "./create_aws_config.sh"]




