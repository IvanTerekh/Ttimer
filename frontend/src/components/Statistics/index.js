import React from 'react';
import AvgCurrent from '../AvgCurrent';
import AvgBest from '../AvgBest';
import ResultList from '../ResultList';
import format from '../../logic/format';
import {avgs} from '../../logic/stats';
import './Statistics.css';

export default function Statistics(props) {
    const containerStyle = {
        width: 9.5
        + (props.stats[100].best !== Infinity)
        + (props.stats[1000].best !== Infinity)
        + format(props.stats.worst).length + 'em'
    };
    const results = props.results.map((result) => result.centis);
    return <div>
        <div className="stats" style={containerStyle}>
            {results.length > 0 &&
            <div className="columns col-gapless">
                <div className="col-6">Avg:<span className="right">{format(props.stats.avg)}</span></div>
                <div className="col-6 right-stats">Count:<span className="right">{props.stats.n}</span></div>
                <div className="col-6">Best:<span className="right">{format(props.stats.best)}</span></div>
                <div className="col-6 right-stats">Worst:<span className="right">{format(props.stats.worst)}</span>
                </div>
                {avgs.map((avgOf) => [
                    <AvgCurrent key={avgOf} avgOf={avgOf} current={props.stats[avgOf].current}/>,
                    <AvgBest key={avgOf} avgOf={avgOf} best={props.stats[avgOf].best}/>,
                ])}
            </div>
            }
        </div>
        <ResultList results={props.results}/>
    </div>
}
