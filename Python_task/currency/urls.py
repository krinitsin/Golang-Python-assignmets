# -*- coding: utf-8 -*-
from django.conf.urls import *
from django.contrib import admin
from django.conf import settings

urlpatterns = [
    url(r'^', include('currency_app.urls')),
    url(r'^admin/', admin.site.urls),
]
