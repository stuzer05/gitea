// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package markdown

import (
	"strings"
	"testing"

	"code.gitea.io/gitea/modules/markup"

	"github.com/stretchr/testify/assert"
)

func TestMathRender(t *testing.T) {
	const nl = "\n"
	testcases := []struct {
		testcase string
		expected string
	}{
		{
			"$a$",
			`<p><code class="language-math is-loading">a</code></p>` + nl,
		},
		{
			"$ a $",
			`<p><code class="language-math is-loading">a</code></p>` + nl,
		},
		{
			"$a$ $b$",
			`<p><code class="language-math is-loading">a</code> <code class="language-math is-loading">b</code></p>` + nl,
		},
		{
			`\(a\) \(b\)`,
			`<p><code class="language-math is-loading">a</code> <code class="language-math is-loading">b</code></p>` + nl,
		},
		{
			`$a$.`,
			`<p><code class="language-math is-loading">a</code>.</p>` + nl,
		},
		{
			`.$a$`,
			`<p>.$a$</p>` + nl,
		},
		{
			`$a a$b b$`,
			`<p>$a a$b b$</p>` + nl,
		},
		{
			`a a$b b`,
			`<p>a a$b b</p>` + nl,
		},
		{
			`a$b $a a$b b$`,
			`<p>a$b $a a$b b$</p>` + nl,
		},
		{
			"a$x$",
			`<p>a$x$</p>` + nl,
		},
		{
			"$x$a",
			`<p>$x$a</p>` + nl,
		},
		{
			"$a$ ($b$) [$c$] {$d$}",
			`<p><code class="language-math is-loading">a</code> (<code class="language-math is-loading">b</code>) [$c$] {$d$}</p>` + nl,
		},
		{
			"$$a$$",
			`<pre class="code-block is-loading"><code class="chroma language-math display">a</code></pre>` + nl,
		},
		{
			"$$a$$ test",
			`<p><code class="language-math display is-loading">a</code> test</p>` + nl,
		},
		{
			"test $$a$$",
			`<p>test <code class="language-math display is-loading">a</code></p>` + nl,
		},
		{
			"foo $x=\\$$ bar",
			`<p>foo <code class="language-math is-loading">x=\$</code> bar</p>` + nl,
		},
	}

	for _, test := range testcases {
		t.Run(test.testcase, func(t *testing.T) {
			res, err := RenderString(markup.NewTestRenderContext(), test.testcase)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, string(res))
		})
	}
}

func TestMathRenderBlockIndent(t *testing.T) {
	testcases := []struct {
		name     string
		testcase string
		expected string
	}{
		{
			"indent-0",
			`
\[
\alpha
\]
`,
			`<pre class="code-block is-loading"><code class="chroma language-math display">
\alpha
</code></pre>
`,
		},
		{
			"indent-1",
			`
 \[
 \alpha
 \]
`,
			`<pre class="code-block is-loading"><code class="chroma language-math display">
\alpha
</code></pre>
`,
		},
		{
			"indent-2",
			`
  \[
  \alpha
  \]
`,
			`<pre class="code-block is-loading"><code class="chroma language-math display">
\alpha
</code></pre>
`,
		},
		{
			"indent-0-oneline",
			`$$ x $$
foo`,
			`<pre class="code-block is-loading"><code class="chroma language-math display"> x </code></pre>
<p>foo</p>
`,
		},
		{
			"indent-3-oneline",
			`   $$ x $$<SPACE>
foo`,
			`<pre class="code-block is-loading"><code class="chroma language-math display"> x </code></pre>
<p>foo</p>
`,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			res, err := RenderString(markup.NewTestRenderContext(), strings.ReplaceAll(test.testcase, "<SPACE>", " "))
			assert.NoError(t, err)
			assert.Equal(t, test.expected, string(res), "unexpected result for test case:\n%s", test.testcase)
		})
	}
}