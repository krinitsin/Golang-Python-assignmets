# -*- coding: utf-8 -*-

from django.conf.urls import *
import currency_app.views

urlpatterns = [
    url(r'^$', currency_app.views.currency_page, name='currency_page'),
    ]
