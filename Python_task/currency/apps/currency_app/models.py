# -*- coding: utf-8 -*-

from django.db import models

class Currency(models.Model):
	NumCode	= models.CharField(max_length=50, verbose_name='NumCode')
	CharCode= models.CharField(max_length=50, verbose_name='CharCode')
	Nominal = models.CharField(max_length=50, verbose_name='Nominal')
	Name = models.CharField(max_length=50, verbose_name='Name')
	Value = models.CharField(max_length=50, verbose_name='Value')

	def __unicode__(self):
		return self.CharCode
