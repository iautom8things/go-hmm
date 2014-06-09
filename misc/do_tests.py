# coding: utf-8

import os, numpy
for i,x in enumerate(numpy.linspace(0.001,0.2,num=20)):
    res = os.system('./hmm -s {} -t {} -o data/trans_prob_{}.dat'.format(i,x,i))

