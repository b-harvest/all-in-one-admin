import { Component } from 'react';
import styled from 'styled-components';
import { size, color } from '@style'
class BaseCard extends Component {
    render() {
        return (
            <Card id={this.props.id} borderColor={this.props.borderColor}>
                <CardTitle>{this.props.title}</CardTitle>
                {this.props.children}
            </Card>
        )
    }
}

const Card = styled.section`
    padding: ${size.base_size_x(5)} ${size.base_size_x(5)} ${size.base_size_x(10)};
    margin: ${size.base_size_x(5)} ${size.base_size_x(3)} ${size.base_size_x(10)};
    background-color: ${color.background};
    border: ${({ borderColor }) => borderColor ? '2px' : '1px'} solid ${({ borderColor }) => borderColor ? borderColor : color.border};
`

const CardTitle = styled.h2`
    margin-bottom: ${size.base_size_x(5)}
`
export default BaseCard