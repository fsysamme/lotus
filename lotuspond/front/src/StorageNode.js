import React from 'react';
import {Cristal} from "react-cristal";

class StorageNode extends React.Component {
  render() {
    return <Cristal
      title={"Storage miner XYZ"}
      initialPosition={'center'}>
      <div className="CristalScroll">
        <div className="StorageNodeInit">
          I'm a node
        </div>
      </div>
    </Cristal>
  }
}

export default StorageNode