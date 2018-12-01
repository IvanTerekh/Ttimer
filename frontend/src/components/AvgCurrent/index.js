import React from 'react';
import '../style.css'
import format from '../../logic/format'

export default function AvgCurrent(props) {
    const avgOf = props.avgOf;

    if (props.current === Infinity) {
        return null;
    } else {
        return <div className="col-6">Avg{avgOf}:<span className="right left-stats">{format(props.current)}</span></div>
    }
}
