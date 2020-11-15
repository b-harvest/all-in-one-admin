import { Component } from 'react';
import styled from 'styled-components';
import { media } from '@style'

class BasicLayout extends Component {

    render() {
        return (
            <Layout>
                {this.props.children}
            </Layout>)
    }
}

const Layout = styled.div`
max-width: 745px;
margin: 0 auto;
text-align: center;
${media.tablet`
    min-width: 722px;
`}
${media.mobile`
    min-width: 375px;
`}

`

export default BasicLayout