#!/usr/bin/env python3

import matplotlib.pyplot as plt
import numpy as np
from scipy.stats import binned_statistic
import matplotlib.ticker as ticker
import math


def set_size(width, fraction=1):
    """ Set aesthetic figure dimensions to avoid scaling in latex.
    Parameters
    ----------
    width: float
            Width in pts
    fraction: float
            Fraction of the width which you wish the figure to occupy
    Returns
    -------
    fig_dim: tuple
            Dimensions of figure in inches
    """
    # Width of figure
    fig_width_pt = width * fraction

    # Convert from pt to inches
    inches_per_pt = 1 / 72.27

    # Golden ratio to set aesthetic figure height
    golden_ratio = (5**.5 - 1) / 2

    # Figure width in inches
    fig_width_in = fig_width_pt * inches_per_pt
    # Figure height in inches
    fig_height_in = fig_width_in * golden_ratio

    fig_dim = (fig_width_in, fig_height_in)

    return fig_dim

def load_compact_dist(fname='compact.txt'):
    """
    input:
    comparing for node_17 - node_17 physical dist 0 approx dist 0
    output:
    python dict of node number and x-coordinates like
    {0: 0.0, 1: 2.0, 2: 4.0}
    """
    out = {}
    with open(fname, 'r') as f:
        for line in f:
            line = line.rstrip()
            splitted = line.split(' ')
            node_src = splitted[2].split('_')[1]
            node_dst = splitted[4].split('_')[1]
            # key = node_src + "_" + node_dst
            val = float(splitted[10])
            out[(node_src, node_dst)] = val

    # print(out)
    return out

def load_latencies(fname='1/pings.txt'):
    """
    input:
    ping node_0 node_1 = 19.314
    output:
    python dict of node number and x-coordinates like
    {0: 0.0, 1: 2.0, 2: 4.0}
    """
    out = {}
    n = 0
    with open(fname, 'r') as f:
        for line in f:
            line = line.rstrip()
            splitted = line.split(' ')
            node_src = splitted[1].split('_')[1]
            node_dst = splitted[2].split('_')[1]
            key = node_src + "_" + node_dst
            val = float(splitted[4])
            out[key] = 2*val
            n += 1
    n = int(math.sqrt(n))

    for i in range(0,n):
        key = str(i) + "_" + str(i)
        out[key] = 0.0

    for i in range(0, n):
        for j in range (0, n):
            key = str(i) + "_" + str(j)
            dist = out[key]
            #for k in range (0, n):
                #key1 = str(i) + "_" + str(k)
                #key2 = str(k) + "_" + str(j)
                #dist1 = out[key1]
                #dist2 = out[key2]
                #if dist > dist1 + dist2:
                    #print("!!!!", str(i) + "_" + str(j) + "=" + str(dist) +"    " + str(i) + "_" + str(k) + "=" + str(dist1) + "    " +str(k) + "_" + str(j) + "=" + str(dist2))


                    #print(out)
    return out

def compute_latencies(fname='coords.txt'):
    """
    input:
    ping node_0 node_1 = 19.314
    output:
    python dict of node number and x-coordinates like
    {0: 0.0, 1: 2.0, 2: 4.0}
    """
    x = {}
    y = {}
    with open(fname, 'r') as f:
        for line in f:
            line = line.rstrip()
            splitted = line.split(' ')
            node = int(splitted[0].split('_')[1])
            x_c = float(splitted[1])
            y_c = float(splitted[2])
            x[node] = x_c
            y[node] = y_c

    return x, y

def load_optime(fname):
    """
    input fname:
    optime-701-SET-node_73 275.875286 163.87528600000002 112
    optime-node_73-node_75 275.875286
    input pairsFile:
    0 node_34 node_84
    optime-701-SET-node_73 275.875286 163.87528600000002 112 0 node_34 node_84
    output:
    {(34, 84): (0.000501547 + 0.061215743)*1000,
     (24, 6) : (0.000547125 + 0.417253898)*1000}
    """
    def parse_one_line(line):
        splitted = line.split(' ')
        latency = float(splitted[1])
        splitted2 = splitted[0].split('-')
        node_num_1 = int(splitted2[1].split('_')[1])
        node_num_2 = int(splitted2[2].split('_')[1])
        return node_num_1, node_num_2, latency

    out = {}
    samePair = {}
    pairs = 0

    lin = []

    with open(fname, 'r') as f:
        flines = f.readlines()
        for line1 in flines:
            lin.append(line1.rstrip())

    # print(lin)

    for l in lin:
        node1, node2, dist = parse_one_line(l)
        out[(node1, node2)] = dist
        pairs += 1
        try:
            samePair[((node1, node2))] += 1
        except:
            samePair[((node1, node2))] = 1

    # print("pairs",pairs,samePair)

    return out


def load_optime_from_sum(fname, isRead):
    """
    input fname:
    optime-701-SET-node_73 275.875286 163.87528600000002 112
    optime-node_73-node_75 275.875286
    input pairsFile:
    0 node_34 node_84
    optime-701-SET-node_73 275.875286 163.87528600000002 112 0 node_34 node_84
    output:
    {(34, 84): (0.000501547 + 0.061215743)*1000,
     (24, 6) : (0.000547125 + 0.417253898)*1000}
    """
    def parse_one_line(line):
        splitted = line.split(' ')
        if len(splitted) == 6:
            if isRead:
                latency = float(splitted[3])
            else:
                latency = float(splitted[1])
            ring = splitted[5]
            return -1,latency,ring
        else:
            ring = splitted[2]
            splitted2 = splitted[0].split('-')
            node_num_1 = int(splitted2[1].split('_')[1])
            node_num_2 = int(splitted2[2].split('_')[1])
        return  node_num_1,  node_num_2, ring

    out = {}
    matchRing = {}

    lin = []

    with open(fname, 'r') as f:
        flines = f.readlines()
        for line1 in flines:
            lin.append(line1.rstrip())

    #print(lin)

    for l in lin:
        out1, out2, ring = parse_one_line(l)
        if out1 == -1:
            matchRing[ring] = out2
        else:
            out[(out1, out2)] = matchRing[ring]
            matchRing = {}

    #print("returning",out)
    return out



def compute_data_points(coord_data, optime_data):
    xdata, ydata = [], []
    for k, latency in optime_data.items():
        node1, node2 = k
        key = str(node1) + "_" + str(node2)
        dist = coord_data[key]
        xdata.append(dist)
        ydata.append(latency)
        #if dist > latency:
        #print("******", node1, node2, dist, latency)
        #if latency > 5 * dist:
        #print("----------", node1, node2, dist, latency)

    #assert dist <=v

    return xdata, ydata


def compute_data_points2(x_c, y_c, optime_data):


    xdata, ydata = [], []
    for k, latency in optime_data.items():
        node1, node2 = k
        dist = 2 * math.sqrt(math.pow((x_c[node2] - x_c[node1]),2) + math.pow((y_c[node2] - y_c[node1]),2))
        xdata.append(dist)
        ydata.append(latency)
        #if dist > latency:
        #print("******", node1, node2, dist, latency)
        #if latency > 5 * dist:
        #print("----------", node1, node2, dist, latency)

    #assert dist <=v

    return xdata, ydata

def scatter_plot(xs, ys, xlabel, ylabel, title, leg=False):
    #width = 505.89
    width = 1000

    fig, ax = plt.subplots(1, 1, figsize=set_size(width, fraction=0.5))
    plt.tight_layout()
    plt.gcf().subplots_adjust(bottom=0.2,left=0.2)

    #fig, ax = plt.subplots()
    ax.scatter(xs, ys, 1.0, alpha=0.5)
    #ax.scatter(xs, ys, 0.2, alpha=0.5)
    #ax.scatter(xs, [y / 2 for y in ys], 1.0, alpha=0.5)
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)
    #ax.set_title(title)

    ax.set_yscale('log')
    ax.set_xscale('log')



    lims = [
        #np.min([ax.get_xlim(), ax.get_ylim()]),  # min of both axes
        np.min([0, 0]),  # min of both axes
        #np.max([ax.get_xlim(), ax.get_ylim()]),  # max of both axes
        np.max([8000, 500]),  # max of both axes
        #np.max([1024, 1024]),  # max of both axes
        #np.max([50, 50]),  # max of both axes
    ]

    # now plot both limits against eachother
    #ax.plot(lims, lims, 'k-', alpha=0.75, zorder=0)
    ax.plot(lims, lims, 'k-', alpha=0.6, label='Ping latency')
    #ax.plot(xs, [ x * 10 for x in xs] , '#A68524', alpha=0.6, label='compact routing latency upper bound')
    #ax.plot(xs, [ x * 18 for x in xs] , 'orange', alpha=0.6, label='ARAs bound K=5 ')
    ax.plot(xs, [ x * 10 for x in xs] , '#A62447', alpha=0.6, label='ARAs bound K=3 ')
    #ax.plot([0, 2], [10, 10], '#3D8BF2', alpha=0.6, label='ARAs latency uper bound')
    '''
    ax.plot([2, 2], [10, 20], '#3D8BF2', alpha=0.6)
    ax.plot([2, 4], [20, 20], '#3D8BF2', alpha=0.6)
    ax.plot([4, 4], [20, 40], '#3D8BF2', alpha=0.6)
    ax.plot([4, 8], [40, 40], '#3D8BF2', alpha=0.6)
    ax.plot([8, 8], [40, 80], '#3D8BF2', alpha=0.6)
    ax.plot([8, 16], [80, 80], '#3D8BF2', alpha=0.6)
    ax.plot([16, 16], [80, 160], '#3D8BF2', alpha=0.6)
    ax.plot([16, 32], [160, 160], '#3D8BF2', alpha=0.6)
    '''


    #ax.plot([32, 64], [320, 320], 'green', lw=2)
    #ax.scatter(xs, [ y * 2 for y in ys] , '1.0')
    #ax.set_aspect('equal')
    #ax.set_xlim(lims)

    #ax.set_ylim([200, 8000.0])
    #ax.set_xlim([32, 256])
    ax.set_ylim([1, 8000.0])
    ax.set_xlim([1, 256])

    #ax.set_xlim([1, 70])

    bins = np.arange(min(xs), max(xs), 2)
    percentilef = lambda v: np.percentile(v, 50, axis=0)
    #print(xs)
    #print(ys)
    #print(bin)
    means = binned_statistic(xs, ys, statistic = percentilef, bins = bins)[0]
    ax.plot(bins[:-1], means, 'go', alpha=0.2, markersize = 4, label="50th percentile")
    percentilef = lambda v: np.percentile(v, 95, axis=0)
    means = binned_statistic(xs, ys, statistic = percentilef, bins = bins)[0]
    ax.plot(bins[:-1], means, 'o', markerfacecolor = "#A62447", markersize = 4, alpha=0.2,label="95th percentile")


    if leg:
        ax.legend(loc='lower right',prop={'size': 6})

    ax.xaxis.set_major_locator(ticker.MultipleLocator(2))
    ax.xaxis.set_major_formatter(ticker.LogFormatter(base=2.0, labelOnlyBase=True, minor_thresholds=None, linthresh=None))

    #plt.grid(True, which="both")



if __name__ == '__main__':

    '''
    xdata, ydata = compute_data_points(load_latencies('../../config/30_random/pings.txt'), load_optime('optime_cruxified.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Experimentally Observed',True)
    #plt.savefig('crux_loc_redis.pdf', format='pdf', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies('../../config/30_random/pings.txt'), load_optime('optime_vanilla.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Vanilla Redis - Experimentally Observed',True)
    #plt.savefig('vanilla_loc_redis.pdf', format='pdf', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies('../../config/30_random/pings.txt'), load_optime('optime_sim.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    #plt.savefig('sim_loc_redis.pdf', format='pdf', dpi=1000)

    #xdata, ydata = compute_data_points(load_latencies('pings_K5.txt'), load_optime('optime_redisK5.txt'))
    #scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    #plt.savefig('sim_loc_redisK5.pdf', format='pdf', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies('pings_N500K3D15.txt'), load_optime('optime_N500K3D15.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    #plt.savefig('sim_loc_redisN500K3D15.pdf', format='pdf', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies('pings_N500K3D150.txt'), load_optime('optime_N500K3D150.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    #plt.savefig('sim_loc_redisN500K3D150.pdf', format='pdf', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies('pings_N500K5D150.txt'), load_optime('optime_N500K5D150.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    #plt.savefig('sim_loc_redisN500K5D150.pdf', format='pdf', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies('pings_N5000K3D15.txt'), load_optime('optime_N5000K3D15.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    plt.savefig('sim_loc_redisN5000K3D15.pdf', format='pdf', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies('pings_N500K3D15.txt'), load_optime('optime_N500K3D15.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    plt.savefig('sim_loc_redisN500K3D15.pdf', format='pdf', dpi=1000)


    xdata, ydata = compute_data_points(load_latencies('pings_N1000K3D120.txt'), load_optime('optime_N1000K3D120.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    plt.savefig('sim_redis_N1000K3D120.pdf', format='pdf', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies('../../config/30_random_D120/pings.txt'), load_optime('optime_N30K3D152.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    plt.savefig('sim_redis_N30K3D152.pdf', format='pdf', dpi=1000)


    xdata, ydata = compute_data_points(load_latencies('../../config/30_random_D120/pings.txt'), load_optime('optime_crux_wide.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    plt.savefig('crux_loc_crdb_wide.pdf', format='pdf', dpi=1000)
    '''

    xdata, ydata = compute_data_points(load_latencies('../data/results/pings.txt'), load_optime('../data/results/min.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Min Cruxified IPFS', True)
    plt.savefig('plot_min.pdf', format='pdf', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies('../data/results/pings.txt'), load_optime('../data/results/vanilla.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Vanilla IPFS', True)
    plt.savefig('plot_vanilla.pdf', format='pdf', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies('../data/results/pings.txt'), load_optime('../data/results/max.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Max Cruxified IPFS', True)
    plt.savefig('plot_max.pdf', format='pdf', dpi=1000)



    """
    xdata, ydata = compute_data_points(load_latencies('../../config/30_random_D120/pings.txt'), load_optime('optime_crux_wide.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    plt.savefig('crux_loc_redis_wide.pdf', format='pdf', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies('../../config/30_random_D120/pings.txt'), load_optime('optime_redis_N30K3D120.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    plt.savefig('crux_sim_redis_wide.pdf', format='pdf', dpi=1000)

    x, y = compute_latencies('coords_N1000K3D120.txt')
    xdata, ydata = compute_data_points2(x, y, load_optime('optime_redis_N1000K3D120.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    plt.savefig('sim_redis_N1000K3D120.pdf', format='pdf', dpi=1000)

    x, y = compute_latencies('coords_N10000K3D120.txt')
    xdata, ydata = compute_data_points2(x, y, load_optime('optime_redis_N10000K3D120.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    plt.savefig('sim_redis_N10000K3D120.pdf', format='pdf', dpi=1000)

    x, y = compute_latencies('coords_N10000K5D120.txt')
    xdata, ydata = compute_data_points2(x, y, load_optime('optime_redis_N10000K5D120.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified Redis - Simulated',True)
    plt.savefig('sim_redis_N10000K5D120.pdf', format='pdf', dpi=1000)
    """

    plt.show()
