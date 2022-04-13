FROM python

ARG UNAME=bob
ARG UID=1000
ARG GID=1000
RUN groupadd -g $GID -o $UNAME
RUN useradd -m -u $UID -g $GID -o -s /bin/bash $UNAME

# Oh no this should be jovyan!
USER bob

# This should be /home/jovyan
WORKDIR /home/bob

ENV PATH=/home/bob/.local/bin:$PATH

RUN pip install jupyterlab

CMD jupyter lab --ip=0.0.0.0
