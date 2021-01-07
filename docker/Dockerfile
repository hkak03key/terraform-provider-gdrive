ARG golang_version=1.15
FROM golang:${golang_version}

ARG tf_version=0.13.0
ARG tf_src_path=/terraform

RUN git clone https://github.com/hashicorp/terraform.git ${tf_src_path} -b v${tf_version} --depth 1
RUN cd ${tf_src_path} && go install

