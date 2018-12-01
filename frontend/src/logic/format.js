export default function format(centis) {
    const hours = Math.floor(centis / (100 * 60 * 60)).toString();
    const minutes = Math.floor(centis % (100 * 60 * 60) / (100 * 60)).toString();
    const seconds = (centis % (100 * 60) / 100).toFixed(2);

    if (hours !== '0') {
        return hours + ':' +
            (minutes.length === 1 ? '0' : '') +
            minutes + ':' +
            (seconds.length === 4 ? '0' : '') +
            seconds;
    } else if (minutes !== '0') {
        return minutes + ':' +
            (seconds.length === 4 ? '0' : '') +
            seconds;
    } else {
        return seconds;
    }
}