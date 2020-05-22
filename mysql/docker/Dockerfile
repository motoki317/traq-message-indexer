FROM mariadb:10.4.13

ENV LANG en_US.utf8

# Install utf8 locale
RUN apt-get update -y && \
	apt-get install -y locales locales-all && \
	rm -rf /var/lib/apt/lists/* && \
	localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8

# Install Mroonga
RUN apt update && \
    apt install mariadb-plugin-mroonga && \
    apt clean && \
    rm -rf /var/lib/apt/lists/* && \
    ln -s /usr/share/mysql/mroonga/install.sql /docker-entrypoint-initdb.d
