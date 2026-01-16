// Tiny JSON tree viewer with collapse/expand and dark-mode awareness
(() => {
  const container = document.getElementById('formattedTree');
  if (!container) return;

  const rawPre = document.getElementById('formattedOutput');
  const toggleBtns = document.querySelectorAll('.view-toggle [data-view]');

  const setView = (view) => {
    toggleBtns.forEach((btn) => {
      btn.classList.toggle('active', btn.getAttribute('data-view') === view);
    });
    if (view === 'raw') {
      container.classList.add('d-none');
      rawPre.classList.remove('d-none');
    } else {
      container.classList.remove('d-none');
      rawPre.classList.add('d-none');
    }
  };

  // Default to showing the raw formatted JSON if there is content, so users immediately see the result text.
  const setInitialView = () => {
    const hasFormatted = rawPre && rawPre.textContent.trim().length > 0;
    setView(hasFormatted ? 'raw' : 'tree');
  };

  toggleBtns.forEach((btn) => {
    btn.addEventListener('click', () => setView(btn.getAttribute('data-view')));
  });

  setInitialView();

  let data;
  try {
    data = JSON.parse(rawPre ? rawPre.textContent : '');
  } catch (e) {
    // Show the raw pane so the user can see the text even if parsing fails.
    setView('raw');
    container.textContent = 'Invalid JSON';
    return;
  }

  const treeEl = document.createElement('div');
  treeEl.className = 'json-tree-root';

  const createNode = (key, value, depth) => {
    const type = Array.isArray(value) ? 'array' : value === null ? 'null' : typeof value;
    const node = document.createElement('div');
    node.className = 'json-node';
    node.style.marginLeft = depth * 16 + 'px';

    const hasChildren = type === 'object' || type === 'array';

    const toggle = document.createElement('button');
    toggle.type = 'button';
    toggle.className = 'json-toggle';
    toggle.setAttribute('aria-label', hasChildren ? 'Toggle' : 'Leaf');
    toggle.textContent = hasChildren ? '▾' : '•';

    const keyEl = document.createElement('span');
    keyEl.className = 'json-key';
    if (key !== null && key !== undefined) keyEl.textContent = key + ':';

    const valueEl = document.createElement('span');
    valueEl.className = 'json-value';

    if (hasChildren) {
      valueEl.textContent = type === 'array' ? `[${value.length}]` : '{ }';
    } else if (type === 'string') {
      valueEl.textContent = ' "' + value + '"';
      valueEl.classList.add('json-string');
    } else {
      valueEl.textContent = ' ' + String(value);
      valueEl.classList.add('json-primitive');
    }

    node.append(toggle, keyEl, valueEl);

    if (hasChildren) {
      const childrenWrapper = document.createElement('div');
      childrenWrapper.className = 'json-children';
      const entries = Array.isArray(value) ? value.entries() : Object.entries(value);
      for (const [childKey, childVal] of entries) {
        childrenWrapper.appendChild(createNode(childKey, childVal, depth + 1));
      }
      node.appendChild(childrenWrapper);

      let collapsed = false;
      const updateToggle = () => {
        toggle.textContent = collapsed ? '▸' : '▾';
        childrenWrapper.style.display = collapsed ? 'none' : 'block';
      };
      const flip = () => {
        collapsed = !collapsed;
        updateToggle();
      };
      toggle.addEventListener('click', () => flip());
      toggle.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          flip();
        }
      });
      node.addEventListener('click', (e) => {
        if (e.target === toggle || e.target.closest('.json-toggle')) return;
        flip();
      });
      childrenWrapper.addEventListener('click', (e) => e.stopPropagation());
      updateToggle();
    }

    return node;
  };

  treeEl.appendChild(createNode(null, data, 0));
  container.innerHTML = '';
  container.appendChild(treeEl);
})();
