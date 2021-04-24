---
template: post
title: What is PyTorch?
shortie: A lightweight introduction to PyTorch.
date: 2017-08-13
categories: deep-learning
tags:
  - python
  - deep learning
---

PyTorch is a scientific computing package for python built by facebook and several other companies and universities. It provides numpy like computation APIs but with strong support
for GPUs. It also supports deep learning specific APIs, for building models and training them. Let's have a quick look at PyTorch.

## Comparing with tensorflow
PyTorch is very much like tensorflow, in a sense that both are used to build complex deep learning models, train them and test them. Both support **computational graphs**, have built-in different kind of layers, optimizers and other utils. But under the hood they are very much different.

Tensorflow, developed by Google was intended to be a library for building production grade models, that can scale better. Hence tensorflow is static in nature. What this means is that, first we define out model from beginning to the end and them execute them. This lets the tensorflow to optimize the graph by, lets say fusing the nodes together or removing redundant nodes. They have even developed a compiler to do all that stuff called **XLA Compiler**. Where as PyTorch is dynamic in nature. Meaning graph gets executed as soon as some node gets attached. The technique used here to achieve this is **Reverse-mode auto-differentiation**.

Other major difference that can be found is that PyTorch is imperative in nature. It means that all the usual python constructs like loop, if-else can be used within PyTorch, unlike tensorflow, where we have to use these constructs as nodes in the graph, like `tf.while_loop`.

PyTorch is designed for rapid prototyping. Hence more and more researchers are now using it in their reseach findings. But tensorflow still remains very popular.

## Getting hands dirty
We import PyTorch with `import torch`. Lets construct a 3x3 matrix. We do something like.
```python
mat = torch.Tensor(3,3)
```
This will immideatily give us a 3x3 matrix, something like
```python
 0.0000e+00  1.0842e-19  9.7721e-38
 3.6902e+19  1.1210e-44 -0.0000e+00
 0.0000e+00  0.0000e+00  0.0000e+00
[torch.FloatTensor of size 3x3]
```

As we discussed above, PyTorch executes the statement immediately instead of defining a session, like we do in tensorflow.

Similarly, there several functions to declare tensors and operate on them. We can also convert them to numpy array and vice versa.

Since we are doing deep learning here, GPUs are must. So in PyTorch, we have this to run our operations in GPU.
```python
if torch.cuda.is_available():
    y = torch.Tensor(5,5)
```

## Going deeper
Ofcourse when doing deep learning, we do backpropagation. And for that we need differentiation. Thankfully PyTorch has a built-in package **autograd** for automatic differentiation.

For this, autograd package has `Variable` class, that wraps a tensor. When these tensors undergo differentiation, there data remians under `.data` attribute and gradient under `.grad` attribute of Variable class.

![Variable]({{site.url}}/assets/Variable.png)
*from pytorch.org*

To do gradient calculation, we call `.backward()`. We also need to speccify `requires_grad=True` when defining variable to calculate gradient of this variable when calling `.backard()`.

Here's some code
```python
x = torch.autograd.Variable(torch.rand(1), requires_grad=True)
y = x + 14
y.backward()
```

We get out gradients accumulated in
```python
x.grad
```
output:
```python
Variable containing:
 1
[torch.FloatTensor of size 1]
```

Well that's pretty much it. In next post we will look at `nn` package and build a CNN for MNIST data using PyTorch.
