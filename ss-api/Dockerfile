FROM public.ecr.aws/lambda/python@sha256:5f59359f025e159ac53ba7bc94d2d724b4d8fc0d5d17375b4b27fb70ae4fa15d as build

RUN dnf install -y unzip && \
    curl -Lo "/tmp/chromedriver-linux64.zip" "https://edgedl.me.gvt1.com/edgedl/chrome/chrome-for-testing/120.0.6099.109/linux64/chromedriver-linux64.zip" && \
    curl -Lo "/tmp/chrome-linux64.zip" "https://edgedl.me.gvt1.com/edgedl/chrome/chrome-for-testing/120.0.6099.109/linux64/chrome-linux64.zip" && \
    unzip /tmp/chromedriver-linux64.zip -d /opt/ && \
    unzip /tmp/chrome-linux64.zip -d /opt/

FROM public.ecr.aws/lambda/python@sha256:5f59359f025e159ac53ba7bc94d2d724b4d8fc0d5d17375b4b27fb70ae4fa15d as prod

RUN dnf install -y atk cups-libs gtk3 libXcomposite alsa-lib \
    libXcursor libXdamage libXext libXi libXrandr libXScrnSaver \
    libXtst pango at-spi2-atk libXt xorg-x11-server-Xvfb \
    xorg-x11-xauth dbus-glib dbus-glib-devel nss mesa-libgbm

# setup pt-BR locale (used for currency formatting)
RUN dnf -y install glibc-locale-source && dnf clean all

RUN localedef -i pt_BR -f UTF-8 pt_BR.UTF-8
ENV LC_ALL pt_BR.UTF-8
ENV LANG pt_BR.UTF-8

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY --from=build /opt/chrome-linux64 /opt/chrome
COPY --from=build /opt/chromedriver-linux64 /opt/

COPY . .

RUN chmod -R 755 /var/task/handlers

CMD [ "main.main" ]
