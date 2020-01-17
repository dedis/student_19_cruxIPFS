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
            #out[key] = 2*val
            out[key] = val
            n += 1
    n = int(math.sqrt(n))

    for i in range(0,n):
        key = str(i) + "_" + str(i)
        out[key] = 0.0

    for i in range(0, n):
        for j in range (0, n):
            key = str(i) + "_" + str(j)
            dist = out[key]

    return out

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

    for l in lin:
        node1, node2, dist = parse_one_line(l)
        out[(node1, node2)] = dist
        pairs += 1
        try:
            samePair[((node1, node2))] += 1
        except:
            samePair[((node1, node2))] = 1

    return out


def compute_data_points(coord_data, optime_data):
    xdata, ydata = [], []
    for k, latency in optime_data.items():
        node1, node2 = k
        key = str(node1) + "_" + str(node2)
        dist = coord_data[key]
        xdata.append(dist)
        ydata.append(latency)

    return xdata, ydata


def scatter_plot_zoom(xs, ys, xlabel, ylabel, title, leg=False):
    width = 505.89

    fig, ax = plt.subplots(1, 1, figsize=set_size(width, fraction=0.5))
    plt.tight_layout()
    plt.gcf().subplots_adjust(bottom=0.2,left=0.2)

    ax.scatter(xs, ys, 1.0, alpha=0.5)
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)

    ax.set_ylim([1, 5000.0])
    ax.set_xlim([1, 400])

    ax.plot(xs, [ x * 1 for x in xs] , '#000000', alpha=0.6, label='Ping latency ')

    bins = np.arange(min(xs), max(xs), 2)
    percentilef = lambda v: np.percentile(v, 50, axis=0)
    means = binned_statistic(xs, ys, statistic = percentilef, bins = bins)[0]
    ax.plot(bins[:-1], means, 'go', alpha=0.2, markersize = 4, label="50th percentile")
    percentilef = lambda v: np.percentile(v, 95, axis=0)
    means = binned_statistic(xs, ys, statistic=percentilef, bins=bins)[0]
    ax.plot(bins[:-1], means, 'o', markerfacecolor="#A62447", markersize=4, alpha=0.2, label="95th percentile")


    if leg:
        ax.legend(loc='lower right', prop={'size': 6})


def scatter_plot(xs, ys, xlabel, ylabel, title, leg=False):
    width = 505.89

    fig, ax = plt.subplots(1, 1, figsize=set_size(width, fraction=0.5))
    plt.tight_layout()
    plt.gcf().subplots_adjust(bottom=0.2,left=0.2)

    ax.scatter(xs, ys, 1.0, alpha=0.5)
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)

    ax.set_yscale('log')
    ax.set_xscale('log')

    lims = [
        np.min([0, 0]),  # min of both axes
        np.max([6000, 500]),  # max of both axes
    ]

    # now plot both limits against eachother
    ax.plot(lims, lims, 'k-', alpha=0.6, label='Ping latency')
    ax.plot(xs, [ x * 10 for x in xs] , '#A62447', alpha=0.6, label='ARAs bound K=3 ')
    ax.set_ylim([1, 6000.0])
    ax.set_xlim([1, 500])

    bins = np.arange(min(xs), max(xs), 2)
    percentilef = lambda v: np.percentile(v, 50, axis=0)
    means = binned_statistic(xs, ys, statistic = percentilef, bins = bins)[0]
    ax.plot(bins[:-1], means, 'go', alpha=0.2, markersize = 4, label="50th percentile")
    percentilef = lambda v: np.percentile(v, 95, axis=0)
    means = binned_statistic(xs, ys, statistic = percentilef, bins = bins)[0]
    ax.plot(bins[:-1], means, 'o', markerfacecolor = "#A62447", markersize = 4, alpha=0.2,label="95th percentile")

    if leg:
        ax.legend(loc='lower right',prop={'size': 6})

    ax.xaxis.set_major_locator(ticker.MultipleLocator(2))
    ax.xaxis.set_major_formatter(ticker.LogFormatter(base=2.0, labelOnlyBase=True, minor_thresholds=None, linthresh=None))


def scatter_plot_read(xs, ys, xlabel, ylabel, title, leg=False):
    width = 505.89

    fig, ax = plt.subplots(1, 1, figsize=set_size(width, fraction=0.5))
    plt.tight_layout()
    plt.gcf().subplots_adjust(bottom=0.2, left=0.2)

    ax.scatter(xs, ys, 1.0, alpha=0.5)
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)

    ax.set_ylim([1, 1500.0])
    ax.set_xlim([1, 400])

    bins = np.arange(min(xs), max(xs), 2)
    percentilef = lambda v: np.percentile(v, 50, axis=0)
    means = binned_statistic(xs, ys, statistic = percentilef, bins = bins)[0]
    ax.plot(bins[:-1], means, 'go', alpha=0.2, markersize = 4, label="50th percentile")
    percentilef = lambda v: np.percentile(v, 95, axis=0)
    means = binned_statistic(xs, ys, statistic = percentilef, bins = bins)[0]
    ax.plot(bins[:-1], means, 'o', markerfacecolor = "#A62447", markersize = 4, alpha=0.2,label="95th percentile")

    if leg:
        ax.legend(loc='lower right',prop={'size': 6})

def scatter_plot_write(xs, ys, xlabel, ylabel, title, leg=False):
    width = 505.89

    fig, ax = plt.subplots(1, 1, figsize=set_size(width, fraction=0.5))
    plt.tight_layout()
    plt.gcf().subplots_adjust(bottom=0.2, left=0.2)

    ax.scatter(xs, ys, 1.0, alpha=0.5)
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)

    ax.set_ylim([1, 4000.0])
    ax.set_xlim([1, 400])

    bins = np.arange(min(xs), max(xs), 2)
    percentilef = lambda v: np.percentile(v, 50, axis=0)
    means = binned_statistic(xs, ys, statistic = percentilef, bins = bins)[0]
    ax.plot(bins[:-1], means, 'go', alpha=0.2, markersize = 4, label="50th percentile")
    percentilef = lambda v: np.percentile(v, 95, axis=0)
    means = binned_statistic(xs, ys, statistic = percentilef, bins = bins)[0]
    ax.plot(bins[:-1], means, 'o', markerfacecolor = "#A62447", markersize = 4, alpha=0.2,label="95th percentile")

    if leg:
        ax.legend(loc='lower right', prop={'size': 6})



def plot_zoom(folder, consistency):
    folder = folder+consistency
    xdata, ydata = compute_data_points(load_latencies(folder+'/data/pings.txt'), load_optime(folder+'/data/min.txt'))
    scatter_plot_zoom(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified IPFS', True)
    plt.savefig(folder+'/graphs/plot_zoom_cruxified_'+consistency+'.png', format='png', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies(folder+'/data/pings.txt'), load_optime(folder+'/data/vanilla.txt'))
    scatter_plot_zoom(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Vanilla IPFS', True)
    plt.savefig(folder+'/graphs/plot_zoom_vanilla_'+consistency+'.png', format='png', dpi=1000)

    plt.show()


def plot_log(folder, consistency):
    folder = folder+consistency
    xdata, ydata = compute_data_points(load_latencies(folder+'/data/pings.txt'), load_optime(folder+'/data/min.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Cruxified IPFS', True)
    plt.savefig(folder+'/graphs/plot_log_cruxified_'+consistency+'.png', format='png', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies(folder+'/data/pings.txt'), load_optime(folder+'/data/vanilla.txt'))
    scatter_plot(xdata, ydata, 'RTT between nodes (ms)', 'W-R pair latency (ms)', 'Vanilla IPFS', True)
    plt.savefig(folder+'/graphs/plot_log_vanilla_'+consistency+'.png', format='png', dpi=1000)

    plt.show()


def plot_read(folder, consistency):
    folder = folder+consistency
    xdata, ydata = compute_data_points(load_latencies(folder+'/data/pings.txt'), load_optime(folder+'/data/read_c.txt'))
    scatter_plot_read(xdata, ydata, 'RTT between nodes (ms)', 'Read latency (ms)', 'Cruxified IPFS', True)
    plt.savefig(folder+'/graphs/plot_read_cruxified_'+consistency+'.png', format='png', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies(folder+'/data/pings.txt'), load_optime(folder+'/data/read_v.txt'))
    scatter_plot_read(xdata, ydata, 'RTT between nodes (ms)', 'Read latency (ms)', 'Vanilla IPFS', True)
    plt.savefig(folder+'/graphs/plot_read_vanilla_'+consistency+'.png', format='png', dpi=1000)

    plt.show()


def plot_write(folder, consistency):
    folder = folder+consistency
    xdata, ydata = compute_data_points(load_latencies(folder+'/data/pings.txt'), load_optime(folder+'/data/write_c.txt'))
    scatter_plot_write(xdata, ydata, 'RTT between nodes (ms)', 'Write latency (ms)', 'Cruxified IPFS', True)
    plt.savefig(folder+'/graphs/plot_write_cruxified_'+consistency+'.png', format='png', dpi=1000)

    xdata, ydata = compute_data_points(load_latencies(folder+'/data/pings.txt'), load_optime(folder+'/data/write_v.txt'))
    scatter_plot_write(xdata, ydata, 'RTT between nodes (ms)', 'Write latency (ms)', 'Vanilla IPFS', True)
    plt.savefig(folder+'/graphs/plot_write_vanilla_'+consistency+'.png', format='png', dpi=1000)

    plt.show()


def plot_all(folder, consistency):
    plot_zoom(folder, consistency)
    plot_log(folder, consistency)
    plot_read(folder, consistency)
    plot_write(folder, consistency)


if __name__ == '__main__':
    folder = 'K3N20D150remoteO2000'
    plot_all(folder, "raft")
    plot_all(folder, "crdt")
