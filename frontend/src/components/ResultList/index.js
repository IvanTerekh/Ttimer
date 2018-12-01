import React from 'react';
import ListItem from '../ListItem'

export default function ResultList(props) {
    const results = props.results;
    return <p className="results">{
        results.map((result) =>
            <ListItem key={result.datetime}
                      value={result.centis}
                      scramble={result.scramble}/>
        )}</p>
}