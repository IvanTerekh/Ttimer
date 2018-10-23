var avgs = [5, 12, 50, 100, 1000];

function avgOf(res, n) {
    if (n > res.length) {
        return Infinity;
    }
    let toIgnore;
    if (n > 2) {
        toIgnore = Math.floor(n / 20) + 1;
    } else {
        toIgnore = 0;
    }
    var latestRes = res.slice(res.length - n);
    ignoreBestWorst(latestRes, toIgnore);
    var sum = latestRes.reduce((a, b) => a + b, 0);
    return Math.round(sum / (n - 2 * toIgnore));
}

function ignoreBestWorst(arr, n) {
    let best = arr.slice(0, n);
    let worst = arr.slice(0, n);
    let bestIndex = Array.from(Array(n).keys());
    let worstIndex = Array.from(Array(n).keys());
    for (let i = n; i < arr.length; i++) {
        if (arr[i] > Math.min.apply(Math, best)) {
            const index = best.indexOf(Math.min.apply(Math, best));
            best[index] = arr[i];
            bestIndex[index] = i;
        }
        if (arr[i] < Math.max.apply(Math, worst)) {
            const index = worst.indexOf(Math.max.apply(Math, worst));
            worst[index] = arr[i];
            worstIndex[index] = i;
        }
    }
    for (let i in bestIndex) {
        arr[bestIndex[i]] = 0;
    }
    for (let i in worstIndex) {
        arr[worstIndex[i]] = 0;
    }
}

function calcStats(res) {
    res = res.map(r => (r.centis));

    let stats = {};
    for (n of avgs) {
        stats[n] = {best: Infinity}
    }

    for (let i = 0; i <= res.length; i++) {
        let slice = res.slice(0, i);
        for (let n of avgs) {
            stats[n].best = Math.min(stats[n].best, avgOf(slice, n));
        }
    }
    stats.best = Math.min(...res);
    stats.worst = Math.max(...res);
    stats.avg = avgOf(res, res.length);
    stats.n = res.length;
    return stats
}

function updateStats(old, res) {
    res = res.map(r => (r.centis));

    let stats = Object.assign({}, old);
    for (let n of avgs) {
        stats[n].best = Math.min(old[n].best, avgOf(res, n));
    }
    stats.best = Math.min(old.best, res[res.length-1]);
    stats.worst = Math.max(old.worst, res[res.length-1]);
    stats.avg = avgOf(res, res.length);
    stats.n = res.length;
    return stats;
}