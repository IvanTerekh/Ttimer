import React from 'react';
import format from '../../logic/format'

export default function ListItem(props) {
    return <span title={props.scramble}>{format(props.value)} </span>;
}
