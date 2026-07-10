FROM python:3.14

ARG UNAME=jovyan
ARG UID=1000
ARG GID=1000
RUN groupadd -g $GID -o $UNAME
RUN useradd -m -u $UID -g $GID -o -s /bin/bash $UNAME

USER jovyan
WORKDIR /home/jovyan

ENV PATH=/home/jovyan/.local/bin:$PATH

RUN pip install jupyterlab==4.6.1 jupyter-server==2.20.0

CMD jupyter lab --ip=0.0.0.0 --ServerApp.base_url="${NB_PREFIX}" --ServerApp.allow_origin="*" --IdentityProvider.token=""
