import 'package:flutter/material.dart';
import 'package:flutter/src/foundation/key.dart';
import 'package:flutter/src/widgets/framework.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';

class CompareSchejDialog extends StatefulWidget {
  final VoidCallback onClose;

  const CompareSchejDialog({Key? key, required this.onClose}) : super(key: key);

  @override
  State<CompareSchejDialog> createState() => _CompareSchejDialogState();
}

class _CompareSchejDialogState extends State<CompareSchejDialog> {
  final _focusNode = FocusNode();

  @override
  initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _focusNode.requestFocus();
    });
  }

  @override
  void dispose() {
    _focusNode.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SchejAppBar(
        isRoot: true,
        underline: false,
        centerTitle: true,
        titleString: 'Compare schej',
        actions: [
          IconButton(onPressed: widget.onClose, icon: const Icon(Icons.close)),
        ],
      ),
      body: Container(
        color: SchejColors.white,
        padding: SchejConstants.pagePadding,
        child: Column(
          children: [
            _buildSearchTextField(),
          ],
        ),
      ),
    );
  }

  Widget _buildSearchTextField() {
    final textField = Padding(
      padding: const EdgeInsets.only(bottom: 20),
      child: TextField(
        focusNode: _focusNode,
        autocorrect: false,
        decoration: const InputDecoration(
          hintText: 'Type a name...',
          prefixIcon: Icon(Icons.search),
        ),
        style: SchejFonts.subtitle,
      ),
    );

    return textField;
  }
}
