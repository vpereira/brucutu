FROM  instrumentisto/dovecot:latest

RUN useradd -ms /bin/bash foo
RUN useradd -ms /bin/bash bar
RUN useradd -ms /bin/bash baz
RUN echo 'foo:thepassword' | chpasswd
RUN echo 'bar:thepassword' | chpasswd
RUN echo 'baz:thepassword' | chpasswd
