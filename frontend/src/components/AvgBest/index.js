import React from 'react';
import '../style.css'
import format from '../../logic/format'

export default function AvgBest(props) {
    const avgOf = props.avgOf;

    if (props.best === Infinity) {
        return null;
    } else {
        return <div className="col-6 right-stats">Best{avgOf}:<span className="right">{format(props.best)}</span></div>
    }
}
