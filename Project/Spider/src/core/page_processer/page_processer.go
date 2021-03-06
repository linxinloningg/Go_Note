// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package page_processer

import (
    "Spider/src/core/common/page"
)

type PageProcesser interface {
    Process(p *page.Page)
    Finish()
}
