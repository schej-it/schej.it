import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter/src/foundation/key.dart';
import 'package:flutter/src/widgets/framework.dart';

class UserAvatar extends StatelessWidget {
  final String src;

  const UserAvatar({Key? key, required this.src}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return CircleAvatar(
      backgroundImage: CachedNetworkImageProvider(
        src,
      ),
    );
  }
}
