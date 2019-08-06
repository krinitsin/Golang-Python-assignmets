# -*- coding: utf-8 -*-
import os
from django.shortcuts import render
from django.views.decorators.csrf import csrf_exempt
from django.shortcuts import redirect
import requests
import xml.etree.ElementTree as ET
from .models import Currency


@csrf_exempt
def currency_page(request):
    if request.method == "POST" or os.path.isfile('./Currencies.xml') != True:
        url = 'https://www.cbr.ru/scripts/XML_daily.asp'
        # creating HTTP response object from given url
        resp = requests.get(url)
        # saving the xml file
        with open('Currencies.xml', 'wb') as f:
            f.write(resp.content)
        return redirect(currency_page)
        # create element tree object
    #preparing data
    tree = ET.parse('Currencies.xml')
    # get root element
    root = tree.getroot()
    # create empty list for currency items
    currencyitems = []
    for currencyVal in root:
        currencyData = Currency()
        currencyData.NumCode = currencyVal[0].text
        currencyData.CharCode = currencyVal[1].text
        currencyData.Nominal = currencyVal[2].text
        currencyData.Name = currencyVal[3].text
        currencyData.Value = currencyVal[4].text
        #saving data to db
        currencyData.save()
        currencyitems.append(currencyData)
    name = root.attrib['name']
    print (name)
    last_update = root.attrib['Date']
    context = {
	'currencyitems':currencyitems,
    "name":name,
    'last_update':last_update
	}

    return render(request, 'currency_page.html', context)
