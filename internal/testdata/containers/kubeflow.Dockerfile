FROM python

ARG UNAME=jovyan
ARG UID=1000
ARG GID=1000
RUN groupadd -g $GID -o $UNAME
RUN useradd -m -u $UID -g $GID -o -s /bin/bash $UNAME

USER jovyan
WORKDIR /home/jovyan

ENV PATH=/home/jovyan/.local/bin:$PATH

RUN pip install jupyterlab

CMD jupyter lab --ip=0.0.0.0
